const m = require("mithril");

const ChannelsList = require("./channels_list");

const MainViewWithChannelsList = {
  view: function(node) {
    return (
        <div class="container-fluid">
          <div class="row">
            <div class="col-sm-3">
              <ChannelsList></ChannelsList>
            </div>
            <div class="col-sm-9">
              {node.children}
            </div>
          </div>
        </div>
    );
  }
};

module.exports = MainViewWithChannelsList;
