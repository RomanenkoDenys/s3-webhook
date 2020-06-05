use Digest::SHA qw(sha256_hex hmac_sha256 hmac_sha256_hex);
use DDP;

my $data = {
    time => '2020-06-05T11:11:47+03:00',
    topic_arn => 'mcs5259999770|myfiles-ash|s3:ObjectCreated:*,s3:ObjectRemoved:*',
    token => 'LfaTdS1ZRSf7thpdQQKTw3GywfjheZ2XCa2njUN1v9UsnuSd',
    url   => 'http://89.208.199.220/webhook',
};


my $step1 = hmac_sha256($data->{time}, $data->{token});
p $step1;
my $step2 = hmac_sha256($data->{topic_arn}, hmac_sha256($data->{time}, $data->{token}));
p $step2;
my $expected_signature = hmac_sha256_hex(
    $data->{url},
    hmac_sha256($data->{topic_arn}, hmac_sha256($data->{time}, $data->{token}))
);

p $expected_signature;
