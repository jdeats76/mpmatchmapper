# mpmatchmapper
Intended to be used as companion web services for the Nakama games server. mpmatchmapper will 1.) create a reduced unique char string given a long char strings 2.) provide a means for clients to pass in the reduced string and fetch the long string. 

Long description:
Nakama server is a game/app server designed to provide social and multi-player gaming backend services.
More information about the Nakama project can be found here
https://github.com/heroiclabs/nakama

The problem mpmatchmapper solves:
Nakama provides real-time multi-player match features with create/invite/join/kick/leave functionality. Part of the process flow for this feature includes a client initiating a match and getting back a matchid in the form of a GUID (e.g. d7725dd8-1645-45fe-8e8b-e96004ea7658 ). To invite players to join the new match, the matchid must be passed along to players chosen to be invited. Nakama does not provide a way to send the matchid to players unless they are already friended though the Nakama friend services. Invited players must copy and paste the full GUID so their game cleint can send a join command.

For clients without the benefit of a full keyboard or the benefit to reference the long string of character (e.g. VR and AR applications) a solution is needed to map the guid to a shorter unique identifier. This is the problem mppatchmapper solves.

The services exposes three end-points. A removemap end-point is provided, but automatic cleanup is also baked in based on time
thresholds which can be configured

writemap:
req:
http://localhost:8087/writemap?key=testkey&matchid=d7725dd8-1645-45fe-8e8b-e96004ea7658
resp:
{"mapid":"37PA"}


readmap:
http://localhost:8087/readmap?key=testkey&mapid=73PA
resp:
{"matchid":"d7725dd8-1645-45fe-8e8b-e96004ea7658"}


removemap:
http://localhost:8087/removemap?key=testkey&mapid=73PA
resp:
{"message":"success"}
