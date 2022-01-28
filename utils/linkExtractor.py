import requests
def linkExtractor(file_path):
    """
    Extracts links from file_path
    returns a list of links

    Also checks if the link is valid
    """
    with open(file_path) as file:
        lines = [line.rstrip('\n') for line in file]
        urls = []
        for line in lines:
            """
            Checking if urls exist
            """
            print("Checking for", line)
            try:
                response = requests.get(line)
                if response.status_code==200:
                    urls.append(line)
            except:
                print("BAD URL, skipping {}".format(line))

        print("VALID URLS")
        for url in urls:
            print(url)
        return urls
            
