var m = require("mithril");

const UserChannels = require("./user_channels");
const getLink = require("./link");

const ChannelsList = {
  view: node => {
    return m(".text-center", [
      getLink(".btn .btn-primary", "/channels/new", "Add new channel"),

      m("hr"),

      m(UserChannels)
    ]);
  }
};

module.exports = ChannelsList;
