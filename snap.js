var system = require('system');
var page = require('webpage').create();
var fs = require('fs');

if (system.args.length !== 4) {
    console.log('Usage: snap.js <some URL> <zoom> <target image name>');
    phantom.exit();
}

function decode_utf8(s) {
  return decodeURIComponent(escape(s));
}

var url = decode_utf8(atob(system.args[1]));
var zoom = system.args[2];
var image_name = system.args[3];
var current_requests = 0;
var last_request_timeout;
var final_timeout;

console.log(url);

page.viewportSize = { width: 1920, height: 1080};
page.settings = { loadImages: true, javascriptEnabled: true };

// If you want to use additional phantomjs commands, place them here
// page.settings.userAgent = 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_8_2) AppleWebKit/537.17 (KHTML, like Gecko) Chrome/28.0.1500.95 Safari/537.17';
page.settings.userAgent = 'Phantom 0.3.0 Diffy Phantom';

// You can place custom headers here, example below.
// page.customHeaders = {
//   'X-Candy-OVERRIDE': 'https://api.live.bbc.co.uk/'
// };

// If you want to set a cookie, just add your details below in the following way.
// phantom.addCookie({
//     'name': 'ckns_policy',
//     'value': '111',
//     'domain': '.bbc.co.uk'
// });

page.onResourceRequested = function(req) {
  current_requests += 1;
};

page.onResourceReceived = function(res) {
  if (res.stage === 'end') {
    current_requests -= 1;
    debounced_render();
  }
};

page.open(url, function(status) {
  if (status !== 'success') {
    console.log('Error with page ' + url);
    phantom.exit();
  } else {
    page.evaluate(function () {
      /* scale the whole body */
      document.body.style.webkitTransform = "scale(2)";
      document.body.style.webkitTransformOrigin = "0% 0%";
      document.body.style.mozTransform = "scale(2)";
      document.body.style.mozTransformOrigin = "0% 0%";
      document.body.style.transform = "scale(2)";
      document.body.style.transformOrigin = "0% 0%";
      /* fix the body width that overflows out of the viewport */
      document.body.style.width = "50%";
    });
  }
});


function debounced_render() {
  clearTimeout(last_request_timeout);
  clearTimeout(final_timeout);

  // If there's no more ongoing resource requests, wait for 3 seconds before
  // rendering, just in case the page kicks off another request
  if (current_requests < 1) {
      clearTimeout(final_timeout);
      last_request_timeout = setTimeout(function() {
          if (current_requests < 1) {
            console.log('Snapping ' + url);
            page.render(image_name);
            phantom.exit();
          } else {
            debounced_render();
          }
      }, 4000);
  }

  // Sometimes, straggling requests never make it back, in which
  // case, timeout after 10 seconds and render the page anyway
  final_timeout = setTimeout(function() {
    console.log('Snapping ' + url);
    page.render(image_name);
    phantom.exit();
  }, 10000);
}
