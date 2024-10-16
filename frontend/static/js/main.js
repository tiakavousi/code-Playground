require.config({ paths: { 'vs': 'https://cdn.jsdelivr.net/npm/monaco-editor@0.31.1/min/vs' }});

require(['vs/editor/editor.main'], function() {
    var editor = monaco.editor.create(document.getElementById('editor'), {
        value: '',
        language: 'python',  // default language
        theme: 'vs-dark'
    });

    // Change language when the user selects a new language
    document.getElementById('language').addEventListener('change', function() {
        var language = this.value;
        monaco.editor.setModelLanguage(editor.getModel(), language);
    });

    // Handle the "Run" button click to execute the code
    document.getElementById('run-btn').addEventListener('click', async function() {
        const language = document.getElementById('language').value;
        const code = editor.getValue();

        const payload = {
            language: language,
            code: code
        };

        try {
            const response = await fetch('http://localhost:8080/execute', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(payload)
            });

            const result = await response.json();
            document.getElementById('output').innerText = result.output || result.error;
        } catch (err) {
            document.getElementById('output').innerText = "Error executing code.";
        }
    });
});

