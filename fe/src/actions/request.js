const m = require("mithril");

function checkAuthAndExtract(xhr, options) {
  if (xhr.status === 401) {
    // User.authState = User.AuthState.SIGNED_OUT;
    m.route.set("/login");
  }
  return {
    response: options.deserialize(xhr.response),
    status: xhr.status
  };
}

module.exports = {
  checkAuthAndExtract: checkAuthAndExtract
};
