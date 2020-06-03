var m = require("mithril");
var User = require("./user");

function checkAuthAndExtract(xhr, options) {
  console.log("check auth");
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
