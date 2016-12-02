#certbot certonly --standalone -d commune.is
certbot renew --hsts
cp /etc/letsencrypt/live/commune.is/cert.pem /etc/letsencrypt/live/commune.is/privkey.pem ./
