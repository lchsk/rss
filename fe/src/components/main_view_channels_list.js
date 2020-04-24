var m = require("mithril");

const ChannelsList = require("./channels_list");

const MainViewWithChannelsList = {
  view: function(node) {
    return m(".container-fluid", [
      m(".row", [
        m(".col-sm-3", m(ChannelsList)),
        m(".col-sm-9", node.children)
      ])
    ]);
  }
};

module.exports = MainViewWithChannelsList;
