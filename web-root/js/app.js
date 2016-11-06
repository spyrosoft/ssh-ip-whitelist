document.getElementById('answer').select();

document.getElementById('answer').addEventListener(
	'keydown',
	submitAnswerOnEnter,
	false
);

function submitAnswerOnEnter(event) {
	clearMessage();
	if (event.keyCode !== 13) { return; } // Return if not Enter key
	var answer = document.getElementById('answer').value;
	if (answer === '') { return; }
	promise.post(
		'/answer',
		{
			'answer' : answer
		}
	).then(answerSubmissionResults);
}

function answerSubmissionResults(error, response, xhr) {
	if (error) {
		printMessage('Something went wrong. The server returned error code: ' + xhr.status);
		return;
	}
	if (response === 'true') {
		printMessage('Success!!');
	} else {
		printMessage('Nope');
	}
}

function printMessage(message) {
	document.getElementById('message').innerHTML = message;
}

function clearMessage() {
	printMessage('');
}