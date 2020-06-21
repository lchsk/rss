const m = require("mithril");

const UserChannels = require("./user_channels");
const getLink = require("./link");

const ChannelsList = {
  view: node => {
    return (
        <div class="text-center">
          {getLink(".btn .btn-primary", "/channels/new", "Add new channel")}
          <hr/>
          <UserChannels></UserChannels>
        </div>
    );
  }
};

module.exports = ChannelsList;
