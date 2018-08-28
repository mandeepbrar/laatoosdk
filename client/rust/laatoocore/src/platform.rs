use request;
use response;
use utils;

type success_callback = fn(response::Response);
type error_callback = fn(i32, String, response::Response);

pub trait Platform {
    fn http_call(&self, request::Request, success_callback, error_callback);
}
