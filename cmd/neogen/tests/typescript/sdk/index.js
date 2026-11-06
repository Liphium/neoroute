// Test double for the SDK imported by generated connectors. It never sends bytes.
// Tests set result to simulate a server response, then inspect payload or trigger handler.
// Share this constructor with the tests so instanceof UserError works normally.
class UserError extends Error {}
class Client {
  constructor(config) {
    this.config = config;
  }

  getConfig() {
    return this.config;
  }

  // Save the event callback so the test can deliver incoming wire values manually.
  receive(name, handler) {
    this.handler = handler;
  }

  // Record the encoded request and return the response arranged by the test.
  async send(name, payload) {
    this.payload = payload;
    return this.result;
  }

  // Response-only routes return the same arranged result without a request payload.
  async sendNoRequest() {
    return this.result;
  }

  // Acknowledgement-only routes record the request and succeed without response data.
  async sendOk(name, payload) {
    this.payload = payload;
  }

  // These variants carry neither request data nor a typed response.
  async sendOkNoRequest() {}

  async sendPing() {}

  // Fire-and-forget routes still let the test inspect the converted request.
  async sendNoResponse(name, payload) {
    this.payload = payload;
  }
}

// Transport setup is a no-op because these tests exercise conversion, not HTTP.
module.exports = { Client, UserError, applyHTTP() {} };
