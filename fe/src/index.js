import m from "mithril";

const SignUpComponent = require('./signup');
const LoginComponent = require('./login');
const LogoutComponent = require('./logout');
const User = require('./user');

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
      m('.row', [
        m('.col-sm', [
          getLink(".btn .btn-block .btn-primary", "/user", "Add new channel"),

          m('hr'),

          m('div', [
            m('div.list-group', [
              getLink(".list-group-item .list-group-item-action", "/user", "Option 1"),
              getLink(".list-group-item .list-group-item-action", "/user", "Option 2"),
              getLink(".list-group-item .list-group-item-action", "/user", "Option 3"),
            ])
          ]),

        ]),
        m('.col-sm-10', "Content"),
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
