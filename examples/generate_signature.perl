use Digest::SHA qw(sha256_hex hmac_sha256 hmac_sha256_hex);
use DDP;

my $data = {
    time => '2019-12-26T19:29:12+03:00',
    topic_arn => 'mcs2883541269|bucketA|s3:ObjectCreated:Put',
    token => 'RPE5UuG94rGgBH6kHXN9FUPugFxj1hs2aUQc99btJp3E49tA',
    url   => 'http://test.com',
};


my $expected_signature = hmac_sha256_hex(
    $data->{url},
    hmac_sha256($data->{topic_arn}, hmac_sha256($data->{time}, $data->{token}))
);

p $expected_signature;
