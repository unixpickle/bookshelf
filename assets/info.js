(function() {

  window.showInfoPopup = function(book) {
    var popup = document.createElement('div');
    popup.className = 'info-popup';

    var image = document.createElement('img');
    image.className = 'cover';
    image.src = book.imageLinks.thumbnail;
    popup.appendChild(image);

    var title = document.createElement('label');
    title.className = 'title';
    title.innerText = book.title;
    popup.appendChild(title);


    var author = document.createElement('label');
    author.className = 'author';
    author.innerText = book.authors ? book.authors.join(', ') : 'No author';
    popup.appendChild(author);

    var description = document.createElement('p');
    description.className = 'description';
    description.innerText = book.description || 'No description available';
    popup.appendChild(description);

    var close = function() {
      document.body.removeChild(shielding);
      document.body.removeChild(popup);
    };

    var closeButton = document.createElement('button');
    closeButton.innerText = 'Close';
    closeButton.addEventListener('click', close);
    popup.appendChild(closeButton);

    var shielding = document.createElement('div');
    shielding.className = 'info-shielding';
    shielding.addEventListener('click', close);

    document.body.appendChild(shielding);
    document.body.appendChild(popup);
  };

})();
