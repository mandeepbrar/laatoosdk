pub enum HttpMethod {
    GET,
    POST,
    PUT,
    DELETE,
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

