var m = require("mithril");

var Footer = {
    view: node => {
        return (
            <footer class="pt-4 my-md-5 pt-md-5 border-top">
                <div class="container">
                    <div className="row">
                        <div className="col-12 col-md text-center">
                            <small className="d-block m-3 text-muted">&copy; 2020</small>
                        </div>
                    </div>
                </div>
            </footer>
        )
    }
};

module.exports = Footer;
