import socket
def create_magic_packet(mac_address):
    """
    根据MAC地址创建魔术包
    """
    # 去除MAC地址中的分隔符
    mac = mac_address.replace(':', '').replace('-', '').replace(' ', '')
    
    # 验证MAC地址长度是否正确（12个字符）
    if len(mac) != 12:
        raise ValueError("Invalid MAC address format")
    
    # 将MAC地址转换为二进制格式
    data = b'FF' * 6 + (mac * 16).encode()
    return bytes.fromhex(data.decode())

def send_wol_packet(mac_address, broadcast_address='255.255.255.255', port=9):
    """
    发送WOL魔术包
    """
    # 创建魔术包
    packet = create_magic_packet(mac_address)
    
    # 创建UDP套接字
    with socket.socket(socket.AF_INET, socket.SOCK_DGRAM) as sock:
        sock.setsockopt(socket.SOL_SOCKET, socket.SO_BROADCAST, 1)  # 启用广播
        sock.sendto(packet, (broadcast_address, port))

#  1-Xc99he7r1uapsklUAC1JZQboJ5XLcHIFm2AQb0j6vCGi6WgwFWLN4M1xeR6XfzrW