use request;
use response;
//use utils;

pub type SuccessCallback = fn(response::Response);
pub type ErrorCallback = fn(response::Response);

pub trait Platform {
    fn http_call(&self, request::HttpRequest, SuccessCallback, ErrorCallback);
}
