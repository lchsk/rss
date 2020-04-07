import m from "mithril";

var Login = {
  current: {},
  submit: () => {
    return m
      .request({
        method: "POST",
        url: "http://localhost:8000/api/authentication",
        data: Login.current,
        withCredentials: true
      })
      .then(result => {
        console.log(result);
      })
      .catch(e => {
        console.log(e.code, e.response, e.message);
      });
  }
};

const UserComponent = {
  oninit: node => {
    m.request({
      method: "GET",
      url: "http://localhost:8000/api/users/1",
      withCredentials: true
    })
      .then(result => {
        console.log(result);
      })
      .catch(e => {
        console.log(e);
      });
  },
  view: node => {}
};

const LoginComponent = {
  view: node => {
    return m(
      "form",
      {
        onsubmit: e => {
          e.preventDefault();
          Login.submit();
        }
      },
      [
        m("input[type=text][placeholder=Email]", {
          oninput: m.withAttr("value", value => {
            Login.current.email = value;
          })
        }),
        m("input[type=password][placeholder=Password]", {
          oninput: m.withAttr("value", value => {
            Login.current.password = value;
          })
        }),
        m("button[type=submit]", "Sign in")
      ]
    );
  }
};

m.route.prefix("#!");

m.route(document.body, "/", {
  "/": {
    render: function() {}
  },
  "/login": {
    render: function() {
      return m(LoginComponent);
    }
  },
  "/user": {
    render: function() {
      return m(UserComponent);
    }
  }
});
