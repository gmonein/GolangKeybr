sed 's/.*value="\(.*\)".*/\1/g'
curl -X POST -F "authenticity_token=iUf2mEf0XsIldpbqADtUD%2FNSpEPxJOMgNUIMEqZ6VXvKPV2jjGPpKjYi5oIZI6LcmU0q0uG9dTzPT0yMKR8rRg=gmonein&user%5Bpassword%5D=Gael26091996%40%2B%2B"

$ curl -X POST -H "Content-Type: application/x-www-form-urlencoded" \
               -d "grant_type=authorization_code&code=362ad374-735c-4f69-aa8e-bf384f8602de&redirect_uri=http://example.com/oauth&client_id=myClientID&client_secret=myClientPassword" \
               https://signin.intra.42.fr/users/sign_in

curl -X POST -d "redirect_uri=google.fr&client_id=08fdc53b395d24655e44ba228b44e2e52b4ed63e683e610a1798bbfa56a39aa6&grant_type=password&username=gmonein&password=Lepassword@42&scope=public" https://api.intra.42.fr/oauth/authorize
