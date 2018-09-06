use std::fmt;

pub enum HttpMethod {
    GET,
    POST,
    PUT,
    DELETE,
}

impl fmt::Display for HttpMethod {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        let printable = match *self {
            HttpMethod::GET => "GET",
            HttpMethod::POST => "POST",
            HttpMethod::PUT => "PUT",
            HttpMethod::DELETE => "DELETE",
        };
        write!(f, "{}", printable)
    }
}

pub struct HttpHeader {
    pub Name: String,
    pub Value: String,
}

pub struct HttpRequest {
    pub URL: String,
    pub Method: HttpMethod,
    pub Headers: Vec<HttpHeader>,
}

pub enum Request<'a> {
    Http(&'a HttpRequest),
}

