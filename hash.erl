% Just a simple implementation of various hashing algorithms so that I have a
% better understanding. It's probably a better idea to just use a C
% wrapper/integration.

-module(hash).
-export([jenkins/2, jenkins_hashing/4]).

jenkins_hashing(Hash, I, Len, _) when I == Len ->
    Hash;
jenkins_hashing(Hash, I, Len, Key) ->
    Hash = Hash + lists:nth(I, Key),
    Hash = Hash + Hash bsl 10,
    Hash = Hash bxor (Hash bsr 6),
    jenkins_hashing(Hash, I + 1, Len, Key).

jenkins(Key, Len) ->
    jenkins_hashing(1, 1, Len, Key).

