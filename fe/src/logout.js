import m from "mithril";

var LogoutComponent = {
  view: node => {
    m.request({
      method: "POST",
      url: "http://localhost:8000/api/logout",
      withCredentials: true
    })
      .then(result => {})
      .catch(e => {});

    m.route.set("/index");
  }
};

export var LogoutComponent;
