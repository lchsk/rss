var m = require("mithril");
const getLink = require("./link");

const NavBar = {
  view: () => (
<nav class="navbar navbar-expand-lg navbar-light bg-light">
  <div class="collapse navbar-collapse" id="navbarToggler">
      <a class="navbar-brand" href="#">RSS</a>
      <div class="mr-auto mt-2 mt-lg-0">

    </div>
    <ul class="navbar-nav mt-2 mt-lg-0">
      <li class="nav-item">
      {getLink(".nav-link", "/login", "Login")}
      </li>
      </ul>
  </div>
</nav>
  )
};

module.exports = NavBar;
