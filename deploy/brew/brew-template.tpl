class Wio < Formula
  desc "An IOT Development Environment"
  homepage "https://github.com/wio/wio"
  url "https://github.com/wio/wio/releases/download/v{{version}}/{{execName}}.{{extension}}"
  version "{{version}}"
  sha256 "{{checksum}}"

  def install
    File.rename("{{execName}}", "wio")
    bin.install "wio"
  end

  test do
    system "#{bin}/wio -v"
  end
end
