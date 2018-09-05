use request;

pub struct HttpService {
    Method: request::HttpMethod,
    URL: String,
}


pub enum Service<'a>  {
    Http(&'a HttpService),
}