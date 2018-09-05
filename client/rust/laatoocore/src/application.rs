//use std::error;
use platform;
use service::{Service};
use request;
use utils::{StringMap};

pub struct Application {
    pf: Box<platform::Platform>
}

impl Application {
    fn execute_service_object(_svc: Service, _service_request: &request::Request, _config: Option<StringMap>) {
        /*var method = get_method(service);
        var req = service_request.get_method_object("http");
        var url = this.getURL(service, req);
        return this.HttpCall(url, method, req.params, req.data, req.headers, config);*/
    }

    fn execute_service(_service_name: String, _service_request: &request::Request, _config: Option<StringMap>) {

    }
}

