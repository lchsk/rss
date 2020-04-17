import m from "mithril";

const SignUpComponent = require('./signup');
const LoginComponent = require('./login');
const LogoutComponent = require('./logout');
const User = require('./user');
const AddNewChannelComponent = require('./add_channel');
var {getErrorMessage, getSingleError} = require("./error");

const UserComponent = {
  oninit: node => {
    User.load();
  },
  view: node => {}
};

function getLink(classes, href, anchor) {
  return m("a" + classes, {
            href: href,
            oncreate: m.route.link,
          }, anchor)
}


var Layout = {
  view: function(node) {

	return m(".container-fluid", [
	  m('.row.top-bar', [
		m(".col-sm-3", [
          m("div", "hello",)
        ]),
		m(".col-sm-6", [
          m("div", "",)
        ]),
		m(".col-sm-3 .text-right", [
          m("div", "user",)
        ]),
	  ]),
      m('.row', [
        m('.col-sm .text-center', [
          getLink(".btn .btn-primary", "/channels/new", "Add new channel"),

          m('hr'),

          m('div', [
            m('div.list-group', [
              getLink(".list-group-item .list-group-item-action", "/user", "Option 1"),
              getLink(".list-group-item .list-group-item-action", "/user", "Option 2"),
              getLink(".list-group-item .list-group-item-action", "/user", "Option 3"),
            ])
          ]),

        ]),
        m('.col-sm-9', node.children),
      ])
    ])
  }
};


var BoxLayout = {
  view: function(node) {

	return m(".text-center", [
      m(".center", [
        m('.top', "Logo"),
        m('', node.children),
        m('.bottom', "Forgot password? SignUP"),
      ]),
	])
  }
};


m.route.prefix("#!");

m.route(document.body, "/", {
  "/": {
    render: function() {
      console.log("route /")
      return m(Layout)
    }
  },
  "/channels/new": {
    render: function() {
      console.log("new channel")
      return m(Layout, m(AddNewChannelComponent))
    }
  },
  "/login": {
    render: function() {
      return m(BoxLayout, m(LoginComponent));
    }
  },
  "/logout": {
    render: function() {
      return m(LogoutComponent);
    }
  },
  "/signup": {
    render: () => {
      return m(SignUpComponent);
    }
  },
  "/user": {
    render: function() {
      return m(UserComponent);
    }
  }
});
