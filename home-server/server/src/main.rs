use tracing_subscriber::{layer::SubscriberExt, util::SubscriberInitExt};

#[tokio::main]
async fn main() {
    let filer_layer = tracing_subscriber::EnvFilter::try_from_default_env().unwrap_or_else(|_| {
        format!(
            "{}=debug,tower_http=debug,axum::rejection=trace",
            env!("CARGO_CRATE_NAME")
        )
        .into()
    });
    let fmt_layer = tracing_subscriber::fmt::layer()
        .compact()
        .with_target(false)
        .with_file(true)
        .with_line_number(true);
    let file_log_layer = tracing_subscriber::fmt::layer()
        .with_writer(tracing_appender::rolling::daily(
            "logs",
            format!("{}.log", env!("CARGO_PKG_NAME")),
        ))
        .with_ansi(false)
        .with_target(true)
        .with_file(true)
        .with_line_number(true);
    tracing_subscriber::registry()
        .with(filer_layer)
        .with(fmt_layer)
        .with(file_log_layer)
        .init();
    tracing::info!("{} v{}", env!("CARGO_PKG_NAME"), env!("CARGO_PKG_VERSION"));
    crate::app::run().await;
}
