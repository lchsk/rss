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

function toggleNavbar() {
  const navbar = document.getElementById("navbarToggler");

  navbar.classList.toggle('collapse');
}

const NavBar = {
  view: () => (
    <nav class="navbar navbar-expand-md navbar-light mb-2">
      <button id="toggler" class="navbar-toggler" type="button" data-toggle="collapse" data-target="navbarToggler"
              aria-controls="navbarToggler" aria-expanded="false" aria-label="Toggle navigation" onclick={toggleNavbar}>
        <span className="navbar-toggler-icon"></span>
      </button>
      <div class="collapse navbar-collapse" id="navbarToggler">
        <a class="navbar-brand" href="/#!/index">
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
