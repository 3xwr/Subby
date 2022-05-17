const searchButton = document.getElementById('search-button');
const searchInput = document.getElementById('search-input');
searchButton.addEventListener('click', () => {
  const inputValue = searchInput.value;
  profile_base_path = "http://localhost:9080/user_profile.html?user="
  window.location.replace(profile_base_path+inputValue);
});