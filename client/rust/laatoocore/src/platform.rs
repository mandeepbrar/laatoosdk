use service;
use http;

pub type SuccessCallback = fn(service::Response);
pub type ErrorCallback = fn(service::Response);

pub trait Platform {
    fn http_call(&self, String, http::HttpMethod, http::HttpRequest, SuccessCallback, ErrorCallback);
}
