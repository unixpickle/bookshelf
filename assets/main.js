window.addEventListener('load', function() {
  if (window.bookInfo.error) {
    var errorDiv = document.createElement('div');
    errorDiv.className = 'error';
    errorDiv.innerText = window.bookInfo.error;
    document.body.appendChild(errorDiv);
    return;
  }

  var booksElement = document.getElementById('books');
  for (var i = 0, len = window.bookInfo.books.length; i < len; ++i) {
    var book = window.bookInfo.books[i];
    var element = document.createElement('li');
    element.className = 'book';

    var image = document.createElement('div');
    image.className = 'cover';
    image.style.backgroundImage = 'url("' + book.imageLinks.thumbnail + '")';
    element.appendChild(image);

    var title = document.createElement('label');
    title.className = 'title';
    title.innerText = book.title;
    element.appendChild(title);

    if (book.authors !== null) {
      var author = document.createElement('label');
      author.className = 'author';
      author.innerText = book.authors.join(', ');
      element.appendChild(author);
    }

    element.addEventListener('click', window.showInfoPopup.bind(null, book));

    booksElement.appendChild(element);
  }
});
