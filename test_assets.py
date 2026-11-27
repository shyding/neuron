import urllib.request
import re

url_base = 'http://localhost:8088/'
with urllib.request.urlopen(url_base) as resp:
    html = resp.read().decode('utf-8', errors='ignore')

links = re.findall(r'href="([^"]+)"', html) + re.findall(r'src="([^"]+)"', html)
print(f"Found {len(links)} asset links in index.html:")
for link in links:
    if link.startswith('http'):
        continue
    url = url_base + link.lstrip('/')
    try:
        req = urllib.request.Request(url)
        with urllib.request.urlopen(req) as r:
            print(f"  {link} -> 200 OK ({len(r.read())} bytes)")
    except urllib.error.HTTPError as e:
        print(f"  {link} -> HTTPError {e.code}")
    except Exception as e:
        print(f"  {link} -> Exception: {e}")
