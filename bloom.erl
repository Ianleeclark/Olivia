-module(bloom_filter).
-export([]).

-record(bloom_filter, {capacity=1000,
                       count=0,
                       bit_array=[]}).

add_key(Filter, _) when Filter#count >= Filter#capacity ->
    % TODO(ian): Change this to have valid errors/error handling.
    1 + 1;
add_key(Filter, Key) ->

    

