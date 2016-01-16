window.addEventListener('load', function() {
  if (window.bookInfo.error) {
    var errorDiv = document.createElement('div');
    errorDiv.className = 'error';
    errorDiv.innerText = window.bookInfo.error;
    document.body.appendChild(errorDiv);
    return;
  }

  var sortedBooks = window.bookInfo.books.slice();
  sortedBooks.sort(function(b1, b2) {
    return b2.updateTime - b1.updateTime;
  });

  var booksElement = document.getElementById('books');
  for (var i = 0, len = sortedBooks.length; i < len; ++i) {
    var book = sortedBooks[i];
    var element = document.createElement('li');
    element.className = 'book';

    var image = document.createElement('div');
    image.className = 'cover';
    image.style.backgroundImage = 'url("' + bookThumbnail(book) + '")';
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

function bookThumbnail(book) {
  if (!book.uploaded) {
    return book.imageLinks.thumbnail + '&usc=0&w=300';
  } else {
    return 'thumbnail?bookId=' + book.id;
  }
}
