var tty = require('tty.js');

var app = tty.createServer({
  shell: 'bash',
  users: {
    foo: 'bar'
  },
  port: 3000
});

app.get('/foo', function(req, res, next) {
  res.send('bar');
});

app.listen();