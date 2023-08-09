# hook-translator

## TODO

### Project Functionality
- add unit tests
    - examing functions/structs to determine what needs to be extracted into a seperate function, or converted to an interface for easier mocking
        - slack client is a good candidate for this, as we'll want to mock it in the handler.
        - investigate how people unit test with Viper, as it'll need to load a config from... somewhere.
- review logging structure to ensure descriptive messages are being returned where necessary.
    - also try to come up with a "requestID" mechanism to track them from start to finish within logs, and include debug info in slack message footer??
- write documentation
- rename project and modules

### 
- add http template type
- add exec template type
