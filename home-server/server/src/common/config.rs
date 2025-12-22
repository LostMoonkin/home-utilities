#[derive(Debug, Clone)]
pub struct Config {
    pub service_host: String,
    pub service_port: String,
}

pub fn load_config() -> Config {
    // Placeholder implementation
    dotenv::dotenv().ok();

    let config = Config {
        service_host: env_get("SERVICE_HOST"),
        service_port: env_parse("SERVICE_PORT"),
    };
    tracing::info!("Loaded config: {:?}", config);
    config
}

#[inline]
fn env_get(key: &str) -> String {
    match std::env::var(key) {
        Ok(v) => v,
        Err(e) => {
            let msg = format!("get env val error: {} {}", key, e);
            tracing::error!(msg);
            panic!("{msg}");
        }
    }
}

#[inline]
fn env_parse<T: std::str::FromStr>(key: &str) -> T {
    let val_str = env_get(key);
    match val_str.parse::<T>() {
        Ok(v) => v,
        Err(_) => {
            let msg = format!("parse env val error: {}, {}", key, val_str);
            tracing::error!(msg);
            panic!("{msg}");
        }
    }
}
