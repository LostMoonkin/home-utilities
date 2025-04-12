import os
import requests
import notify
import wol
import time
import datetime

source_uri = "http://192.168.1.3/api/v2.0"
source_api_key = os.environ.get("source_nas_api_key")
dest_uri = "http://192.168.1.13/api/v2.0"
dest_api_key = os.environ.get("dest_nas_api_key")
dest_mac = "2A:14:6C:11:0F:C8"
replication_task_name_list = ["BIWIN_NV7400_bk", "KIOXIA-G3_bk"]
wol_first_check_delay = 30
wol_check_interval = 5
wol_timeout = 600
monitor_replication_task_interval = 30
monitor_replication_task_timeout = 21600 # 6 Hours


def check_system_ready(host_uri, api_key):
    system_ready_url = host_uri + "/system/ready"
    r = requests.get(url=system_ready_url, headers={
        "Authorization": "Bearer " + api_key,
        "ContentType": "application/json"
    }, timeout=2)
    return r.status_code == 200 and r.json() == True

def get_replication_tasks(host_uri, api_key, task_name_list = []):
    url = host_uri + "/replication"
    r = requests.get(url, headers={
        "Authorization": "Bearer " + api_key,
        "ContentType": "application/json"
    }, timeout=2)
    if r.status_code != 200:
        print(f"Get replication task error, status_code={r.status_code}, body={r.content}")
        return {
            "code": 4,
            "msg": f"Get replication task error, status_code={r.status_code}"
        }
    body = r.json()
    res = []
    for task in body:
        if task["name"] in task_name_list:
            res.append(task)
            print("Find replication task: id={0}, name={1}".format(task["id"], task["name"]))
    return {
        "code": 0,
        "data": res
    }

def get_replication_task(host_uri, api_key, task_id):
    url = f"{host_uri}/replication/id/{task_id}"
    r = requests.get(url, headers={
        "Authorization": "Bearer " + api_key,
        "ContentType": "application/json"
    }, timeout=2)
    if r.status_code == 200:
        return r.json()
    return None

def trigger_replication_task(host_uri, api_key, task_id):
    url = f"{host_uri}/replication/id/{task_id}/run"
    r = requests.post(url, data={}, headers={
        "Authorization": "Bearer " + api_key,
        "ContentType": "application/json"
    }, timeout=10)
    if r.status_code != 200:
        print(f"Trigger replication task failed, task_id={task_id} code={r.status_code}, body={r.content}")
        return
    print(f"Trigger replication task success, task_id={task_id}")

def shutdown_system(host_uri, api_key):
    url = host_uri + "/system/shutdown"
    r = requests.post(url, data={}, headers={
        "Authorization": "Bearer " + api_key,
        "ContentType": "application/json"
    }, timeout=2)
    return r.status_code == 200

def run_backup():
    # Check source system ready
    if not check_system_ready(source_uri, source_api_key):
        return {
            "code": 1,
            "msg": "source system not ready, task aborted."
        }
    # Get Replication tasks
    task_resp = get_replication_tasks(source_uri, source_api_key, replication_task_name_list)
    if task_resp["code"] != 0:
        return task_resp
    replication_tasks = task_resp["data"]
    # Send WOL magic packet
    try:
        wol.send_wol_packet(dest_mac)
        print(f"WOL packet sent to {dest_mac}, start dealy for {wol_first_check_delay} seconds.")
    except Exception as e:
        print(f"Error: {str(e)}")
        return {
            "code": 2,
            "msg": f"WOL failed: {str(e)}"
        }
    time.sleep(wol_first_check_delay)
    cnt = 0
    # Waiting for dest system ready
    while not check_system_ready(dest_uri, dest_api_key):
        if cnt * wol_check_interval >= wol_timeout:
            return {
                "code": 3,
                "msg": "WOL timeout."
            }
        print(f"Destination system boot not ready. Sleep for {wol_check_interval} second and try again.")
        cnt += 1
        time.sleep(wol_check_interval)
    # Trigger replication task
    triggered_task_id_list = []
    for task in replication_tasks:
        if "state" in task:
            state = task["state"]
            if "state" in state and state["state"] in ["RUNNING", "PENDING"]:
                print("Replication task is running or pending, skiped. name={}".format(task["name"]))
                continue
        print("Trigger replication task: name={0}, id={1}".format(task["id"], task["name"]))
        triggered_task_id_list.append(task["id"])
        trigger_replication_task(source_uri, source_api_key, task["id"])
    # Monitor task status
    time.sleep(monitor_replication_task_interval)
    task_res = []
    finished_task_id_list = []
    cnt = 0
    while len(finished_task_id_list) < len(triggered_task_id_list) and cnt * monitor_replication_task_interval < monitor_replication_task_timeout:
        for id in triggered_task_id_list:
            task = get_replication_task(source_uri, source_api_key, id)
            if task is not None and task["id"] not in finished_task_id_list:
                if "state" in task:
                    state = task["state"]
                    if state is not None and state["state"] in ["ERROR", "FINISHED"]:
                        finished_task_id_list.append(task["id"])
                        state_time = datetime.datetime.now()
                        if "datetime" in state and "$date" in state["datetime"]:
                            state_time = datetime.datetime.fromtimestamp(state["datetime"]["$date"] / 1000)
                        print("Task {0}: state={1}, time={2}".format(task["name"], state["state"], state_time.strftime("%FT%XZ")))
                        task_res.append("Task {0}: state={1}, time={2}".format(task["name"], state["state"], state_time.strftime("%FT%XZ")))
        cnt += 1
        print(f"Monitor replication tasks, total task count= {len(triggered_task_id_list)}, finished task count = {len(finished_task_id_list)}")
        time.sleep(monitor_replication_task_interval)
    # Timeout.
    if len(finished_task_id_list) != len(triggered_task_id_list):
        return {
            "code": 4,
            "title": "检测备份任务超时！！！",
            "msg": "\n".join(task_res)
        }
    # All tasks finished, shutdown destitation system.
    r = shutdown_system(dest_uri, dest_api_key)
    if not r:
        return {
            "code": 5,
            "title": "目标主机关机失败！！！",
            "msg": "\n".join(task_res)
        }
    return {
        "code": 0,
        "msg": "\n".join(task_res)
    }

if __name__ == '__main__':
    res = run_backup()
    if res["code"] != 0:
        notify.send(res["title"] if "title" in res else "Truenas备份失败!!!","code: {0}, msg: {1}".format(res["code"], res["msg"]))
    else:
        notify.send("Truenas备份成功", res["msg"])