var system = require('system'),
    page = require('webpage').create();

page.viewportSize = { width: 1920, height: 1080 };
page.open(system.args[1], function (status) {
  if (status !== 'success') {
    console.log('Unable to access the network!');
    phantom.exit();
  } else {
    setTimeout(function() {
      page.render(system.args[2]);
      phantom.exit();
    }, system.args[3]); // wait for sign to render and fonts/images to load and settle
  }
});
