const path = require('path');
const notifier = require('node-notifier');
const express = require('express');
const app = express();
const bodyParser = require('body-parser');
const port = process.env.PORT || 9000;

app.use(bodyParser.json());
app.get('/health', (req, res) => res.status(200).send());
app.post('/notify', (req, res) => {
    notify(req.body, reply => res.send(reply))
});

app.listen(port, () => console.log(`server is up and running on port: ${port}`));

const notify = ({title, message}, cb) => {
    notifier.notify(
        {
            title: title || 'Unknown title',
            message: message || 'Unknown message',
            icon: path.join(__dirname, 'gophertuts.png'),
            sound: true,
            wait: true,
            reply: true,
            closeLabel: 'Completed?',
            timeout: 15,
        },
        (err, response, reply) => {
            cb(reply)
        }
    );
};
