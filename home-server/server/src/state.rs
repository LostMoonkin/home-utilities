pub struct AppState {
    // Add shared state fields here
    pub config: crate::common::config::Config,
}

impl AppState {
    pub fn new(config: crate::common::config::Config) -> Self {
        AppState { config }
    }
}
