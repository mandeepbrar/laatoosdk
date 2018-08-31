pub enum HttpMethod {
    GET,
    POST,
    PUT,
    DELETE,
};

pub struct HttpHeader {
    Name String,
    Value String
}

pub struct HttpRequest {
    URL: String,
    Method: HttpMethod,
    Headers: []HttpHeader
};

