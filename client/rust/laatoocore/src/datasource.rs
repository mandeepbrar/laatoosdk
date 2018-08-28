use request;

fn execute_service_object(service: &service::Service, service_request: &request::Request, config: Option<StringMap>) {
    var method = get_method(service);
    var req = service_request.get_method_object("http");
    var url = this.getURL(service, req);
    return this.HttpCall(url, method, req.params, req.data, req.headers, config);
}

fn execute_service(service_name: String, service_request: &request::Request, config: Option<StringMap>) {

}



  HttpCall(url, method, params, data, headers, config=null) {
    let service = this;
    var promise = new Promise(
      function (resolve, reject) {
        if (method === "" || url === "") {
          reject(service.buildHttpSvcResponse(Response.InternalError, 'Could not build request', url));
          return;
        }
        let successCallback = function(response) {
          if (response.status < 300) {
            let res = service.buildHttpSvcResponse(Response.Success, "", response);
            resolve(res);
          } else {
            reject(service.buildHttpSvcResponse(Response.Failure, "", response));
          }
        };
        let errorCallback = function(response) {
          reject(service.buildHttpSvcResponse(Response.Failure, "", response));
        };
        if(method == 'DELETE' || method == 'GET') {
          data = null;
        }
        if(!headers) {
          headers = {}
        }
        headers[Application.Security.AuthToken] = Storage.auth;
        let req = {
          method: method,
          url: url,
          data: data,
          headers: headers,
          params: params,
          responseType: 'json'
        };
        if(config) {
          req = Object.assign({}, req, config)
        }
        console.log("Request.. ",req);
        axios(req).then(successCallback, errorCallback);
      });
    return promise;
  }

  createFullUrl(url, params) {
    if (params != null && Object.keys(params).length != 0) {
      return url + "?" + Object.keys(data).map(function(key) {
        return [key, data[key]].map(encodeURIComponent).join("=");
      }).join("&");
    }
    return url
  }

  buildHttpSvcResponse(code, msg, res) {
    if(res instanceof Error) {
      return this.buildSvcResponse(code, msg, res, {});
    }
    return this.buildSvcResponse(code, msg, res.data, res.headers, res.status);
  }

  buildSvcResponse(code, msg, data, info, statuscode) {
    var response = {};
    switch (code) {
      default:
        response.code = code;
        response.message = msg;
        response.data = data;
        response.info = info;
        response.statuscode = statuscode;
    }
    console.log(response);
    return response;
  }

