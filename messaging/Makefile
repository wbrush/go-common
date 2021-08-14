all: mockgen

mockgen:
	#TODO add all packages and interfaces that need to be mocked
	mockgen -destination ./mocks/mock_message.go bitbucket.org/optiisolutions/go-common/messaging Message
	mockgen -destination ./mocks/mock_pub_sub.go bitbucket.org/optiisolutions/go-common/messaging PublisherSubscriber