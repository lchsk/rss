const m = require("mithril");
const getLink = require("./link");
const LoginComponent = require("./login");

function loginInfo() {
  if (LoginComponent.signedIn === true) {
    return (
        <p class="nav-link disabled mb-0">{"Hello, " + LoginComponent.username}</p>
    )
  }

  return getLink(".nav-link .mb-0", "/login", "Login");
}
function logoutInfo() {
  if (LoginComponent.signedIn === true) {
    return getLink(".nav-link .mb-0", "/logout", "Logout");
  }

  return '';
}

const NavBar = {
  view: () => (
    <nav class="navbar navbar-expand-lg navbar-light mb-2">
      <div class="collapse navbar-collapse" id="navbarToggler">
        <a class="navbar-brand" href="#">
          <img src="data/text2011.png" alt="rss"/>
        </a>
        <div class="mr-auto mt-2 mt-lg-0"></div>
        <ul class="navbar-nav mt-0 mb-0 mt-lg-0">
          <li class="nav-item mb-0">{loginInfo()}</li>
          <li class="nav-item mb-0">{logoutInfo()}</li>
        </ul>
      </div>
    </nav>
  )
};

module.exports = NavBar;
