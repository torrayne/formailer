package main

const emailTemplate = `<html lang="en">
<head>
    <style>
        body {
            font-family: sans-serif;
        }

        #wrapper {
            padding: 30px 20px;
            border: 1px solid #ccc;
            border-radius: 4px;
            margin: 20px 0 10px;
        }

        .content {
            width: 100%;
            max-width: 800px;
            margin: 0;
            box-sizing: border-box;
        }

        h1 {
            font-size: 1.8rem;
            color: #212121;
            font-weight: bold;
            margin: 0 0 10px;
        }

        table {
            border-collapse: collapse;
        }

        th,
        td {
            vertical-align: top;
            padding: 10px 10px;
            border-top: 1px solid #ddd;
            text-align: left;
            color: #404040;
            font-size: 1.1rem;
            font-weight: 400;
        }

        th {
            font-weight: bold;
        }

        .attribute {
            padding: 0 20px;
            text-decoration: none;
            color: #404040;
        }
    </style>
</head>
<body>
    <div id="wrapper" class="content">{{ . }}</div>
    <p class="content"><a class="attribute" href="https://atwood.io">Powered by Formailer © Atwood.io</a></p>
</body>
</html>`
