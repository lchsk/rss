var m = require("mithril");

function getLink(classes, href, anchor) {
  return m(
    "a" + classes,
    {
      href: href,
      oncreate: m.route.link
    },
    anchor
  );
}

module.exports = getLink;
