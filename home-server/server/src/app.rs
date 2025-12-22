use crate::common;
pub async fn run() {
    let config = common::config::load_config();
    let state = crate::state::AppState::new(config);
}
