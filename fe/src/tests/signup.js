var mq = require("mithril-query")
var o = require("ospec")

var SignUpComponent = require("../signup.js")
var User = require("../user.js")

o.spec("SignUpComponent", function() {
  o("form is rendered when signed out", function() {
    User.authState = User.AuthState.SIGNED_OUT;
    User.load = function() {}

    var out = mq(SignUpComponent)

    out.should.have("form");
    out.should.have("div#signup-error");
    out.should.have("input[type=text][placeholder=Email]");
    out.should.have("input[type=text][placeholder=Username]");
    out.should.have("input[type=password][placeholder=Password]");
    out.should.have("input[type=password][placeholder='Repeat password']");
    out.should.have("button[type=submit]");
  })
  o("loading div is rendered when in default state", function() {
    var out = mq(SignUpComponent)

    out.should.have("div");
  })
})
