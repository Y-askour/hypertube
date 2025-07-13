defmodule HypertubeTest do
  use ExUnit.Case
  doctest Hypertube

  test "greets the world" do
    assert Hypertube.hello() == :world
  end
end
