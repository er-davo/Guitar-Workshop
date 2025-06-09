from service import server
import os

if __name__ == "__main__":
    server.run_server(os.getenv('PORT'))